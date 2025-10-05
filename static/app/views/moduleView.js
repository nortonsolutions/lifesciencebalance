/**
 * The module view is the main view for the module page.
 * 
 * Renders and provides supporting JavaScript for module page,
 * including ModuleElement and Timer react components.
 *
 * localStorage will contain entries that look like this:
 * 
 * { 
 *    userIdmoduleId: {
 *       this.userAnswers: {},
 *       timePassed: Number, (seconds),
 *       currentElementIndex: Number
 *    } 
 * }
 * 
 */

// import { TimerContextProvider } from './moduleProvider.js';
import { Timer } from './moduleTimer.js';
 
class ModuleView {
    constructor(router, eventDepot) {

        this.router = router;
        this.eventDepot = eventDepot;
        this.msg = new SpeechSynthesisUtterance();
        this.markedOptions = {
            breaks: true
        }
        this.getTemplates();
    }

    async getTemplates() {
        this.template = await getTemplate(`app/views/templates/module.hbs`);
        this.elementTemplate = await getTemplate('app/views/templates/element.hbs');
    }

    async render(moduleId, dom = document.getElementById("main")) {
  
        await this.preLoad("/cdn/js/marked/marked.min.js");
        await this.preLoad("/cdn/js/ajax/libs/react/18.2/react.production.min.js")
        await this.preLoad("/cdn/js/ajax/libs/react/18.2/react-dom.production.min.js")
        // await this.preLoad("/cdn/js/babel-standalone_7.18.8/babel.min.js")
        
        // await this.preLoadMathJax();
        // await this.preLoad("/cdn/js/MathJax-Master/MathJax.js?config=TeX-AMS-MML_HTMLorMML");

        this.currentElementIndex = 0;
        this.currentElementID = '';
        this.userAnswers = {};
        this.running = true;

        this.timePassed = 0;

        this.moduleId = moduleId;
        this.module = JSON.parse(await handleGet(`/module/${moduleId}`));

        let elementIdList = await handleGet(`/module/${moduleId}/moduleelement`);
        this.module.elementIds = (JSON.parse(elementIdList)).map(el => el.element_id);
        this.module.totalElements = Number(this.module.elementIds.length);

        dom.innerHTML = await this.template(this.module);

        this.userId = localStorage.getItem('userId');
        this.courseId = localStorage.getItem('courseId');

        this.timeLimit = this.module.time_limit==0? Infinity : Number(this.module.time_limit)*60; //seconds

        // These will have to be passed in
        this.reviewMode = false;
        this.showAnswers = false;

        // Is a module already in progress for this user?
        if (localStorage.getItem(this.userId + this.moduleId)) {
            let moduleObj = JSON.parse(localStorage.getItem(this.userId + this.moduleId));
            this.userAnswers = moduleObj.userAnswers;
            this.timePassed = moduleObj.timePassed;
            this.currentElementIndex = moduleObj.currentElementIndex;
        } 

        this.renderTimer();

        this.addButtonListeners();
        this.renderElement();

        this.addSpeechSynthesis();

        // MathJax:
        this.Preview = {
            CreatePreview: function () {
                MathJax.Hub.Queue(
                    ["Typeset",MathJax.Hub,document.getElementById("elementMarked")],
                    ["Typeset",MathJax.Hub,document.getElementById("text0")],
                    ["Typeset",MathJax.Hub,document.getElementById("text1")],
                    ["Typeset",MathJax.Hub,document.getElementById("text2")],
                    ["Typeset",MathJax.Hub,document.getElementById("text3")],
                    ["PreviewDone",this]
                );
            },
            PreviewDone: function () {
                }
        };
                
        // this.Preview.callback = MathJax.Callback(["CreatePreview",this.Preview]);
        // this.Preview.callback.autoReset = true;  // make sure it can run more than once

    }

    tick = () => {

        if (this.timePassed >= this.timeLimit) {
            this.running = false;
        }

        if (this.running) {
            this.timePassed += 1;
            this.eventDepot.fire('tick', this.timePassed)
        }
    }

    // React Timer feature (No JSX!  Pure JS: https://babeljs.io/repl)
    renderTimer() {
        if (!this.reviewMode) {

            const root = ReactDOM.createRoot(document.getElementById('moduleTimer'));
            root.render(
                React.createElement(Timer,{ eventDepot: this.eventDepot, timeLimit: this.timeLimit},null)
            );

            if (this.running) {
                let myVar = setInterval(() => {
                    if (this.running) {
                        this.tick();
                    } else {
                        // TODO: submitQuiz();
                        clearInterval(myVar);
                    }
                }, 1000)
            }
        }
    }

    renderElement = async () => {
        
        // Save the current answers, time passed, and current element index
        localStorage.setItem(this.userId+this.moduleId, JSON.stringify({
            userAnswers: this.userAnswers,
            timePassed: this.timePassed,
            currentElementIndex: this.currentElementIndex
        }));

        // console.log(this.currentElementID,this.module.elementIds,this.currentElementIndex)
        this.currentElementID = this.module.elementIds[this.currentElementIndex];
        
        let elementRaw = await handleGet('/element/' + this.currentElementID)
        this.currentElement = JSON.parse(elementRaw);
        

        this.currentElement.reviewMode = this.reviewMode;
        this.currentElement.showAnswers = this.showAnswers;
        this.currentElement.totalElements = this.module.totalElements;
        this.currentElement.currentElementNumber = this.currentElementIndex + 1;
        
        let elementDom = document.getElementById('element');
        elementDom.innerHTML = this.elementTemplate(this.currentElement);

        // Render the markdown text
        var elementMarked = marked(document.getElementById('elementRaw').innerText, this.markedOptions);
        document.getElementById('elementMarked').innerHTML = elementMarked;

        // Speech synthesis stuff
        this.msg.text = document.getElementById('elementMarked').innerText;
        let answers = document.getElementById('responseCheckboxes').querySelectorAll('textarea');
        if (answers) {
            this.msg.text += Array.from(answers)
                .filter(answer => answer.value.length > 0)
                .map((answer,index) => (index+1) + answer.value)
                .join();
        }
            
        this.setButtonVisibility();
        this.applyCheckboxHandlers();

        if (!this.reviewMode && this.currentElement.type == "projectSubmission") {
            this.addProjectSubmissionHandling();
        }

        this.applyImageLink();
        
        this.populateCurrentAnswer();

        // this.Preview.callback();

    }
    
    hideElementInterface = () => {
        document.getElementById('moduleControls').classList.add('d-none');
        document.getElementById('element').classList.add('d-none');
        document.getElementById('moduleTimer').classList.add('d-none');
    }
    
    setButtonVisibility = () => {
        let btnPreviousElement = document.getElementById('previousElement');
        let btnNextElement = document.getElementById('nextElement');
        let btnFinishModule = document.getElementById('finishModule');
    
        // Visibility
        if (this.currentElementIndex == 0) {
            btnPreviousElement?.classList.add('d-none');
            btnFinishModule?.classList.add('d-none');
            btnNextElement?.classList.remove('d-none');
        } else if (this.currentElementIndex + 1 == this.module.totalElements){
            btnNextElement?.classList.add('d-none');
            btnPreviousElement?.classList.remove('d-none');
            if (!this.reviewMode) btnFinishModule?.classList.remove('d-none');
        } else {
            btnFinishModule?.classList.add('d-none');
            btnNextElement?.classList.remove('d-none');
            btnPreviousElement?.classList.remove('d-none');
        }
    
        if (document.getElementById('elementType').value == "projectSubmission") {
            btnFinishModule?.classList.add('d-none');
        }
    
    }
    
    addButtonListeners = () => {
        let btnPreviousElement = document.getElementById('previousElement');
        let btnNextElement = document.getElementById('nextElement');
        let btnFinishModule = document.getElementById('finishModule');
    
        // Click events
        btnPreviousElement.addEventListener('click', (e) => {
            this.saveCurrentAnswer();
            this.currentElementIndex--;
            this.renderElement();
        })
    
        btnNextElement.addEventListener('click', (e) => {
            this.saveCurrentAnswer();
            this.currentElementIndex++;
            this.renderElement();
        })
    
        btnFinishModule.addEventListener('click', (e) => {
            this.saveCurrentAnswer();
            this.submitModule();
        })
    }
    
    applyCheckboxHandlers = () => {
        let checkboxArray = Array.from(document.querySelectorAll(".checkbox"));
        checkboxArray.forEach(checkbox => {
            checkbox.addEventListener('change', (e) => {
                let currentItem = e.target;
                if (document.getElementById('elementType').value == 'single' && currentItem.checked) {
                    // Unselect the others
                    checkboxArray.forEach(item => {
                        if (item.name != currentItem.name) item.checked = false;
                    })
                }
            })
        })
    }
    
    // TODO: implement and test this
    addProjectSubmissionHandling = () => {
        Array.from(document.querySelectorAll('.custom-file-input')).forEach(el => {
            el.addEventListener('change', e => {
                let filename = String(e.target.files[0].name);
                let ending = Array.from(filename).splice(-3,3).join('');
                let firstpart = filename.substring(0,filename.indexOf('.')).slice(0,8);
                e.target.labels[0].innerText = firstpart + "..." + ending;
            })
        })
    
        document.getElementById('projectSubmissionForm').addEventListener('submit', e => {
    
            // This is like a hybrid cross of save and submit, for projects
            if (! this.userAnswers[elementId]) {
                this.userAnswers[elementId] = {};
            }
    
            this.userAnswers[elementId].projectFile = e.target.elements[0].files[0].name;// filename;
            
            var formData = new FormData(e.target);
            formData.append('userAnswers', JSON.stringify(this.userAnswers));
            // formData.append('timePassed', store? store?.getState()["timePassed"] : 0);
            formData.append('projectFile', e.target.elements[0].files[0].name);
    
    
            handleFormPost('/module/projectSubmission/' + this.courseId + '/' + this.moduleId, formData, (response) => {
                
                let returnToCoursePage = "<a href='/app.html?page=course&detail=" + this.courseId + "'>Return to course page</a>"; 
                document.getElementById('feedback').innerHTML = JSON.parse(response).feedback + "..." + returnToCoursePage;
                localStorage.removeItem(this.userId + this.moduleId);
                hideElementInterface();
            }) 
    
            e.preventDefault();
        })
    
    }
    
    applyImageLink = () => {
        document.querySelector('.elementImage')?.addEventListener('click', (e) => {
            window.open('/public/uploads/' + e.target.id,'CourseApp image',"top=500,left=500,width=800,height=800"); 
    
        })
    }

    populateCurrentAnswer = () => {

        let elementType = document.getElementById('elementType').value;
        if (this.userAnswers[this.currentElementID]) {
    
            if (this.userAnswers[this.currentElementID].answer) {
                this.userAnswers[this.currentElementID].answer.forEach((check, index) => {
                    document.getElementById('checkbox' + index).checked = check;
                });
                document.getElementById('answerText').value = this.userAnswers[this.currentElementID].answerText;
                document.getElementById('answerEssay').value = this.userAnswers[this.currentElementID].answerEssay;
            }
            
            if (this.userAnswers[this.currentElementID].projectFile) {
                document.getElementById('userProjectLink').innerHTML="Existing project: <a href='/public/projects/" + this.userAnswers[this.currentElementID].projectFile + "'>" + this.userAnswers[this.currentElementID].projectFile + "</a>"
            }
    
            if (this.reviewMode && elementType != "instructional" && elementType != "projectSubmission") {
                let correctOrIncorrect = document.getElementById('correctOrIncorrect');
                if (this.userAnswers[this.currentElementID].correct && this.userAnswers[this.currentElementID].correct == true) {
                    correctOrIncorrect.style.color = 'green';
                    correctOrIncorrect.style.fontStyle = 'bold';
                    correctOrIncorrect.innerText = 'Correct!';
                } else {
                    correctOrIncorrect.style.color = 'red';
                    correctOrIncorrect.style.fontStyle = 'bold';
                    correctOrIncorrect.innerText = 'Incorrect';
                }
            }
    
        }
    }
    
    saveCurrentAnswer = () => {
    
        
        if (!this.reviewMode) {
            if (! this.userAnswers[this.currentElementID]) {
                this.userAnswers[this.currentElementID] = {};
            }
        
            this.userAnswers[this.currentElementID].answer = [
                    document.getElementById('checkbox0').checked,
                    document.getElementById('checkbox1').checked,
                    document.getElementById('checkbox2').checked,
                    document.getElementById('checkbox3').checked
               ];
            this.userAnswers[this.currentElementID].answerText = document.getElementById('answerText').value;
            this.userAnswers[this.currentElementID].answerEssay = document.getElementById('answerEssay').value;
            
            // Save this.userAnswers - moved to renderElement
            
        }
    }
    
    submitModule = () => {
        handlePost('/module/grade/' + this.courseId + '/' + this.moduleId, {
            userAnswers: this.userAnswers, 
            timePassed: this.timePassed
        }, (response) => {
    
            let returnToCoursePage = "<a href='/app.html?page=course&detail=" + this.courseId + "'>Return to course page</a>"; 
            document.getElementById('feedback').innerHTML = JSON.parse(response).feedback + " ... " + returnToCoursePage;
            localStorage.removeItem(this.userId + this.moduleId);
            this.hideElementInterface();
        })
    }

    /* Thanks to Wes Bos for inspiring this stuff */
    addSpeechSynthesis = () => {

        let voices = [];
        const voicesDropdown = document.querySelector('[name="voice"]');
        const options = document.querySelectorAll('[type="range"], [name="text"]');
        const speakButton = document.querySelector('#speak');

        const populatevoices = () => {
            voices = speechSynthesis.getVoices();
            // console.log(voices);

            voicesDropdown.innerHTML = voices
                .filter(voice => voice.name.match(/English/))
                .map(voice => `<option value="${voice.name}">${voice.name} (${voice.lang})</option>`)
                .join('');
        }

        const setVoice = (e) => {
            this.msg.voice = voices.find(voice => voice.name == e.target.value);
        }

        const toggle = (startOver = true) => {
            speechSynthesis.cancel();

            if (startOver) {
                speechSynthesis.speak(this.msg);
            }
        }

        const setOption = (e) => {
            this.msg[e.target.name] = e.target.value;
            toggle();
        }

        speechSynthesis.addEventListener('voiceschanged', populatevoices);
        voicesDropdown.addEventListener('change', setVoice);
        options.forEach(option => option.addEventListener('change', setOption));
        speakButton.addEventListener('click', toggle);
        
        populatevoices();
    }

    preLoad = (scripting) => {
        return new Promise((resolve, reject) => {
            var script = document.createElement("script");
            script.src = scripting;
            script.addEventListener("load", resolve);
            script.addEventListener("error", reject);
            document.body.appendChild(script);
        });
    };

    preLoadMathJax = (scripting) => {
        return new Promise((resolve, reject) => {
            var script = document.createElement("script");
        //   script.src = scripting;
            script.type="text/x-mathjax-config";
            script.innerHTML=
            'MathJax.Hub.Config({showProcessingMessages: false,tex2jax: { inlineMath: [["$","$"],["\\(","\\)"]] },"HTML-CSS": {availableFonts:[], preferredFont: null, webFont: null});';
            script.addEventListener("load", resolve);
            script.addEventListener("error", reject);
            document.body.appendChild(script);
        });
    };

}
 
 export { ModuleView };
 