/* Norton's rendition of the MVC Property Mgmt Page (June 2019) */
class PropertyDetailsView {
  constructor(eventDepot) {
    this.eventDepot = eventDepot;
    this.currentPropertyId = -1;
    this.templateRetrieved = false;

    fetch("app/views/templates/propertyDetails.hbs")
      .then(response => response.text())
      .then(text => {
        this.template = Handlebars.compile(text);
        this.templateRetrieved = true;
      });
  }
}

PropertyDetailsView.prototype.clearSlate = async function() {
  let context = { id: this.currentPropertyId };

  while (!this.templateRetrieved) await delay();
  document.getElementById("propertyDetailPlaceholder").innerHTML = this.template(context);
  document.getElementById("title").focus();
  this.addButtonListeners();
};

PropertyDetailsView.prototype.hideInputForm = async function() {
  while (!this.templateRetrieved) await delay();
  document.getElementById("inputForm").hidden = true;
  document.getElementById("btnSave").hidden = true;
  document.getElementById("btnCancel").hidden = true;
};

PropertyDetailsView.prototype.showInputForm = function() {
  document.getElementById("inputForm").hidden = false;
  document.getElementById("btnSave").hidden = false;
  document.getElementById("btnCancel").hidden = false;
};

PropertyDetailsView.prototype.getCurrentPropertyId = function() {
  return Number(document.getElementById("uniqueId").innerText);
};

PropertyDetailsView.prototype.setCurrentPropertyId = function(propertyId) {
  this.currentPropertyId = propertyId;
  document.getElementById("uniqueId").innerText = Number(propertyId);
};

PropertyDetailsView.prototype.getFormInputValues = function() {
  let property = {};

  for (let el of document.getElementById("propertyDetail").elements) {
    switch (el.name) {
      case "id":
        property[el.name] = Number(el.value);
        break;

      case "gender":
        if (el.checked) property[el.name] = el.value;
        break;

      default:
        property[el.name] = el.value;
        break;
    }
  }

  return property;
};

PropertyDetailsView.prototype.populatePropertyDetails = async function(currentProperty) {
  while (!this.templateRetrieved) await delay();
  document.getElementById("propertyDetailPlaceholder").innerHTML = this.template(currentProperty);
  this.addButtonListeners();
  document.title = "Edit " + currentProperty.name;
};

PropertyDetailsView.prototype.validateForm = function() {
  let form = document.getElementById("propertyDetail");

  if (!form.checkValidity()) {
    form.reportValidity();
    return false;
  } else return true;
};

PropertyDetailsView.prototype.addButtonListeners = function() {
  document
    .getElementById("propertyDetail")
    .addEventListener("submit", e => { this.eventDepot.fire('saveProperty', e); });

  document
    .getElementById("btnCancel")
    .addEventListener("click", e => { this.eventDepot.fire('cancel', e); });
  document.getElementById("title").focus();
};

export { PropertyDetailsView };
