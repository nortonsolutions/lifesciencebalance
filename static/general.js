const getOptions = () => {
  return {
    method: "GET",
  };
};

// Assumes Handlebars has been preloaded
const getTemplate = async (url) => {
  try {
    let response = await fetch(url, { ...getOptions(), cache: "no-cache" });
    let text = await response.text();
    let template = Handlebars.compile(text);
    return template;
  } catch (error) {
    throw error;
  }
};

const handleGet = async (url) => {
  try {
    let response = await fetch(url, { ...getOptions(), cache: "no-cache" });
    let text = await response.text();
    return text;
  } catch (error) {
    throw error;
  }
};

const getParamsFromRequest = (request) => {
  let { parameters } = request;

  let paramObj = {};

  if (request.parameters[0]) {
    for (let el of parameters) {
      let [name, value] = el.split("=");

      if (value) {
        value = value.replace(/\+/g, " ");
        paramObj[name] = value;
      } else {
        paramObj[name] = null;
      }
    }
  }
  return paramObj;
};

// Borrowed from D3 code
let shuffle = function (array, i0, i1) {
  if ((m = arguments.length) < 3) {
    i1 = array.length;
    if (m < 2) i0 = 0;
  }
  var m = i1 - i0,
    t,
    i;
  while (m) {
    i = (Math.random() * m--) | 0;
    (t = array[m + i0]), (array[m + i0] = array[i + i0]), (array[i + i0] = t);
  }
  return array;
};

/* Register Handlebars helpers for views */
const registerHandlebarsHelpers = () => {
  Handlebars.registerHelper("iconize", function (data) {
    let str;

    switch (data) {
      case "male":
        str = "<i class='fas fa-male'></i>";
        break;
      case "female":
        str = "<i class='fas fa-female'></i>";
        break;
      default:
        str = "?";
    }
    return new Handlebars.SafeString(str);
  });

  Handlebars.registerHelper("booleanCheckboxHelper", function (boolean) {
    let str = "";
    if (boolean) str = "checked";
    return new Handlebars.SafeString(str);
  });

  Handlebars.registerHelper("plusOne", function (value) {
    return value + 1;
  });

  Handlebars.registerHelper("limitLength", function (string) {
    if (string.length > 22) {
      return string.substring(0, 22) + "...";
    } else return string;
  });

  Handlebars.registerHelper(
    "hideButton",
    function (button, currentElementNumber, totalElements) {
      switch (button) {
        case "previous":
          if (currentElementNumber == 1) return "d-none";
          break;
        case "next":
          if (currentElementNumber == totalElements) return "d-none";
          break;
        case "done":
          if (currentElementNumber != totalElements) return "d-none";
          break;
      }
      return "";
    }
  );

  Handlebars.registerHelper("disabledIf", function (boolean) {
    let str = "";
    if (boolean) str = "disabled";
    return new Handlebars.SafeString(str);
  });

  Handlebars.registerHelper("hideIf", function (boolean) {
    let str = "";
    if (boolean) str = "d-none";
    return new Handlebars.SafeString(str);
  });

  Handlebars.registerHelper("hideIfNot", function (boolean) {
    let str = "";
    if (!boolean) str = "d-none";
    return new Handlebars.SafeString(str);
  });

  Handlebars.registerHelper(
    "hideIfNotEqual",
    function (string1, string2, string3) {
      let str = "";
      if (string1 != string2) {
        if (string3) {
          if (string1 != string3) str = "d-none";
        } else {
          str = "d-none";
        }
      }
      return new Handlebars.SafeString(str);
    }
  );

  Handlebars.registerHelper(
    "hideIfEqual",
    function (string1, string2, string3) {
      let str = "";
      if (string1 == string2) {
        if (string3) {
          if (string1 == string3) str = "d-none";
        } else {
          str = "d-none";
        }
      }
      return new Handlebars.SafeString(str);
    }
  );

  Handlebars.registerHelper(
    "selectedIfEqual",
    function (string1, string2, string3) {
      let str = "";
      if (string1 == string2 || string1 == string3) str = "selected";
      return new Handlebars.SafeString(str);
    }
  );

  Handlebars.registerHelper("shortDate", function (date) {
    return date.toLocaleDateString();
  });

  Handlebars.registerHelper("jsDate", function (date) {
    var formatter = new Intl.DateTimeFormat("en-us", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
    });

    let [{ value: month }, , { value: day }, , { value: year }] =
      formatter.formatToParts(date);
    return year + "-" + month + "-" + day;
  });

  Handlebars.registerHelper("calculateRows", function (string) {
    if (!string || string.length == 0) return 4;
    return Math.ceil(string.length / 45);
  });

  Handlebars.registerHelper("calculateRowsShort", function (string) {
    if (!string || string.length == 0) return 2;
    return Math.ceil(string.length / 37);
  });

  Handlebars.registerHelper("correctOrIncorrrect", function (correct) {
    let str = "";
    if (correct) {
      str = '<strong style="color: green">Correct!</strong>';
    } else {
      str = '<strong style="color: red">Incorrect</strong>';
    }

    return new Handlebars.SafeString(str);
  });

  Handlebars.registerHelper("selectedForRole", function (roles, role) {
    if (roles && roles.includes(role)) {
      return "selected";
    } else return "";
  });

  Handlebars.registerHelper("divide", function (one, two) {
    return (one / two) * 100;
  });

  Handlebars.registerHelper("reportedHelper", function (reported) {
    return reported ? "reported" : "";
  });

  Handlebars.registerHelper("secondsToMMSS", function (seconds) {
    return (
      Math.floor(seconds / 60) +
      ":" +
      (seconds % 60).toLocaleString("en-US", { minimumIntegerDigits: 2 })
    );
  });

  Handlebars.registerHelper("homeContentHelper", function (homeContent) {
    return new Handlebars.SafeString(homeContent);
  });
};
