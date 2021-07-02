function deleteRow(id) {
  var row = document.getElementById(id);
  row.remove();
}

function addRow() {
  var row = document.createElement("tr");
  var cols = [];
  var table = document.getElementById("table-body");
  var rows = table.rows.length;

  row.id = "deletable-{x}".replace("{x}", rows + 1);

  for (let i = 1; i <= 8; i++) {
    cols.push(document.createElement("td"));
  }

  for (c in cols) {
    row.appendChild(cols[c]);
  }

  cols[0].innerHTML = '<input type="checkbox"  name="include{x}" checked=true>';
  cols[0].innerHTML = cols[0].innerHTML.replace("{x}", rows + 1);
  
  cols[1].innerHTML = '<input type="text" name="name-{x}" value="' + document.getElementById("newColName").value + '">';
  cols[1].innerHTML = cols[1].innerHTML.replace("{x}", rows + 1);
  
  cols[2].innerHTML = '<input type="text" name="type-{x}" value="' + document.getElementById("newColType").value  + '">';
  cols[2].innerHTML = cols[2].innerHTML.replace("{x}", rows + 1);
  
  cols[3].innerHTML = '<td><input type="checkbox" name="unique-{x}"></td>'
  cols[3].innerHTML = cols[3].innerHTML.replace("{x}", rows + 1);

  
  cols[4].innerHTML =
    '<select name="regex-opt{x}">' +
    document.getElementById("template").value
    "</select>";
  cols[4].innerHTML = cols[4].innerHTML.replace("{x}", rows + 1);

  cols[5].innerHTML = '<input type="checkbox"  name="custom-{x}">';
  cols[5].innerHTML = cols[5].innerHTML.replace("{x}", rows + 1);
  
  cols[6].innerHTML = '<input type="text"      name="regex-{x}">';
  cols[6].innerHTML = cols[6].innerHTML.replace("{x}", rows + 1);
  
  cols[7].innerHTML =
    '<button type="button" onclick=\'deleteRow("{id}")\'>' +
    "  Delete Row" +
    "</button>";
  cols[7].innerHTML = cols[7].innerHTML.replace("{x}", rows + 1);
  cols[7].innerHTML = cols[7].innerHTML.replace("{id}", "deletable-{x}".replace("{x}", rows + 1));

  table.appendChild(row);
}
