/**
 * HomeView
 * @author Norton 2022
 */
class HomeView {

    constructor(router, eventDepot) {
        this.router = router;
        this.eventDepot = eventDepot;
    }

    async render(context) {
        // if (!context) var context = {};
        context.indexActive = 'active';

        this.template = await getTemplate('app/views/templates/home.hbs');
        this.dom = document.getElementById('main');
        this.dom.innerHTML = this.template(context);

        this.router.setRouteLinks(this.dom);
        document.title = `nortonApp`;
        window.scrollTo(0,0);
        this.addEventListeners();
        this.renderCourses(context)
    }

    addEventListeners = () => {}

    renderCourses = async (user) => {
        let courses = JSON.parse(await handleGet(`/user/${user.id}/course`));
        this.template = await getTemplate('app/views/templates/myCourses.hbs');
        this.coursesDom = document.getElementById('myCourses');
        this.coursesDom.innerHTML = this.template(courses);
        this.router.setRouteLinks(this.coursesDom);
    }
}

export { HomeView };