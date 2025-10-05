// TODO: simplify and consolidate?  Some overlap here.
swalAlert = (e) => {
  // catches runtime errors, e.g. throw new Error(....);
  let timerInterval;
  Swal.fire({
    title: e.error,
    html: e.message + "<br> <b></b>",
    timer: 3000,
    icon: "warning",
    timerProgressBar: true,
    didOpen: () => {
      Swal.showLoading();
      const b = Swal.getHtmlContainer().querySelector("b");
      timerInterval = setInterval(() => {
        b.textContent = Swal.getTimerLeft();
      }, 100);
    },
    willClose: () => {
      clearInterval(timerInterval);
    },
  }).then((result) => {
    if (result.dismiss === Swal.DismissReason.timer) {
      window.location.href = "/index.html";
    }
  });
};

postOptions = (data) => {
  return {
    method: "POST",
    body: JSON.stringify(data),
    headers: { "Content-Type": "application/json; charset=utf-8" },
  };
};

handlePost = (url, data, callback) => {
  fetch(url, postOptions(data))
    .then((response) => response.json())
    .catch((error) => {
      callback(JSON.stringify(error.message));
    })
    .then((response) => {
      callback(JSON.stringify(response));
    });
};

putOptions = (data) => {
  return {
    method: "PUT",
    body: JSON.stringify(data),
    headers: { "Content-Type": "application/json; charset=utf-8" },
  };
};

handlePut = (url, data, callback) => {
  fetch(url, putOptions(data))
    .then((response) => response.json())
    .catch((error) => {
      callback(JSON.stringify(error.message));
    })
    .then((response) => {
      callback(JSON.stringify(response));
    });
};

handlePutTextResponse = (url, data, callback) => {
  fetch(url, putOptions(data))
    .then((response) => response.text())
    .catch((error) => {
      callback(error.message);
    })
    .then((response) => {
      callback(response);
    });
};

formPostOptions = (data) => {
  return {
    method: "POST",
    body: data,
  };
};

handleFormPost = (url, data, callback) => {
  fetch(url, formPostOptions(data))
    .then((response) => response.json())
    .catch((error) => {
      callback(JSON.stringify(error.message));
    })
    .then((response) => {
      callback(JSON.stringify(response));
    });
};

handlePostTextResponse = (url, data, callback) => {
  fetch(url, postOptions(data))
    .then((response) => response.text())
    .catch((error) => {
      callback(error.message);
    })
    .then((response) => {
      callback(response);
    });
};

deleteOptions = (data) => {
  return {
    method: "DELETE",
    body: JSON.stringify(data),
    headers: { "Content-Type": "application/json; charset=utf-8" },
  };
};

handleDelete = (url, data, callback) => {
  fetch(url, deleteOptions(data))
    .then((response) => response.text())
    .catch((error) => {
      callback(error.message);
    })
    .then((response) => {
      callback(response);
    });
};

plotAxes = (ctx, maxX, maxY) => {
  let yZero = canvas.height - 40;
  let xZero = 30
  let xUnit = (canvas.width - xZero) / maxX;
  let yUnit = (canvas.height - 80) / maxY;

  ctx.clearRect(0, 0, canvas.width, canvas.height);
  ctx.beginPath();
  ctx.moveTo(xZero, 30);
  ctx.lineTo(xZero, yZero);
  ctx.lineTo(xZero + canvas.width, yZero);
  ctx.strokeStyle = "black";
  ctx.stroke();
  ctx.closePath();

  // add x-axis labels every ten units from 0 to maxX
  ctx.font = "12px Arial";
  ctx.fillStyle = "black";
  for (let i = 0; i <= maxX; i+=10) {
    ctx.fillText(i, xZero + i * xUnit, yZero + 20);
  }
  // add y-axis labels every ten units from 0 to maxY
  for (let i = 0; i <= maxY; i+=10) {
    ctx.fillText(i, xZero - 20, yZero - i * yUnit);
  }

  // add x-axis label
  ctx.fillText("Time", xZero + canvas.width / 2, yZero + 40);
  // add y-axis label
  ctx.save();
}

plotArray = (ctx, array, title, color, labelOffset, maxY) => {
  let yZero = canvas.height - 40;
  let x = xZero = 30

  let xUnit = (canvas.width - xZero) / array.length;
  let yUnit = (canvas.height - 80) / maxY;

  ctx.font = "12px Arial";
  ctx.fillStyle = color;
  ctx.fillText(title, x + labelOffset, 20);

  ctx.beginPath();
  // Each element in the array is a y-coordinate on the graph.
  // For each y-coordinate, draw a line from the previous y-coordinate.
  // Each x-coordinate is the index of the array, starting at 0, and increasing by canvasWidth/array.length each time.
  ctx.moveTo(xZero, yZero);
  for (let i = 0; i < array.length; i++) {
    x += xUnit;
    y = yZero - array[i] * yUnit;
    ctx.lineTo(x, y);

    // every 10th point, place a label on the graph
    if (i % 10 == 0) {
      ctx.fillText(array[i].toFixed(1), x - 7, y + 14);
    }
    
  }
  ctx.strokeStyle = color;
  ctx.stroke();
  ctx.closePath();

}

getTableFromArray = (items, {columnNames}) => {
  let table = document.createElement('table');
  let header = table.createTHead();
  columnNames.forEach((columnName) => {
    let th = document.createElement('th');
    th.style.width = "100px";
    th.innerText = columnName;
    header.appendChild(th);
  });

  items.forEach((item) => {
    let row = table.insertRow();
    for (let i = 0; i < columnNames.length; i++) {
      let cell = row.insertCell();
      cell.innerText = item[columnNames[i]];
    }
  });

  return table;
}