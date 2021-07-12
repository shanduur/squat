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

  for (let i = 1; i <= 11; i++) {
    cols.push(document.createElement("td"));
  }

  var newColName = document.getElementById("newColName")
  var name = newColName.value;
  var type = document.getElementById("newColType").value;

  var i = 0;
  addNo(cols[i++], name, rows+1);
  addInclude(cols[i++], name);
  addName(cols[i++], name);
  addType(cols[i++], name, type);
  addLength(cols[i++], name, 1);
  addPrecision(cols[i++], name, 0);
  // addUnique(cols[i++], name);
  addNullable(cols[i++], name);
  addData(cols[i++], name);
  addCustom(cols[i++], name);
  addREGEX(cols[i++], name);
  addDelete(cols[i++], rows+1);

  for (c in cols) {
    row.appendChild(cols[c]);
  }

  table.appendChild(row);

  if (/^(.*\d+)$/.test(name)) {
    newColName.value = name.replace(/\d+/, String(rows+1));
  } else {
    newColName.value = name + "_" + String(rows+1);
  }

}

{/* <td><input type="number" name="order-{{ .Name }}" value="{{ .Order }}" style="width:48px;"></td> */}
function addNo(td, name, order) {
  var input = document.createElement("input");

  input.setAttribute("type", "number");
  input.setAttribute("name", "order-{{ .Name }}".replace("{{ .Name }}", name));
  input.setAttribute("value", String(order));
  input.setAttribute("style", "width:48px;");

  td.appendChild(input);
}

{/* <td><input type="checkbox" name="include-{{ .Name }}" checked=true></td> */}
function addInclude(td, name) {
  var input = document.createElement("input");

  input.setAttribute("type", "checkbox");
  input.setAttribute("name", "include-{{ .Name }}".replace("{{ .Name }}", name));
  input.setAttribute("checked", true);

  td.appendChild(input);
}

{/* <td><input type="text" name="name-{{ .Name }}" value="{{ .Name }}"></td> */}
function addName(td, name) {
  var input = document.createElement("input");

  input.setAttribute("type","text");
  input.setAttribute("name", "name-{{ .Name }}".replace("{{ .Name }}", name));
  input.setAttribute("value", "{{ .Name }}".replace("{{ .Name }}", name));

  td.appendChild(input);
}

{/* <td><input type="text" name="type-{{ .Name }}" value="{{ .Type }}"></td> */}
function addType(td, name, type) {
  var input = document.createElement("input");

  input.setAttribute("type", "text");
  input.setAttribute("name", "type-{{ .Name }}".replace("{{ .Name }}", name));
  input.setAttribute("value", "{{ .Type }}".replace("{{ .Type }}", type));


  td.appendChild(input);
}

{/* <td><input type="number" name="length-{{ .Name }}" value="{{ .Length }}"></td> */}
function addLength(td, name, length) {
  var input = document.createElement("input");

  input.setAttribute("type", "number");
  input.setAttribute("name","length-{{ .Name }}".replace("{{ .Name }}", name));
  input.setAttribute("name","length-{{ .Name }}".replace("{{ .Name }}", name));
  input.setAttribute("value", String(length));

  td.appendChild(input);
}

{/* <td><input type="number" name="precision-{{ .Name }}" value="{{ .Precision }}"></td> */}
function addPrecision(td, name, precision) {
  var input = document.createElement("input");

  input.setAttribute("type", "number");
  input.setAttribute("name", "precision-{{ .Name }}".replace("{{ .Name }}", name));
  input.setAttribute("value", String(precision));

  td.appendChild(input);
}

{/* <td><input type="checkbox" name="nullable-{{ .Name }}" readonly="readonly"></td> */}
function addCustom(td, name) {
  var input = document.createElement("input");

  input.setAttribute("type", "checkbox");
  input.setAttribute("name", "nullable-{{ .Name }}".replace("{{ .Name }}", name));

  td.appendChild(input);
}

{/* <td><input type="checkbox" name="unique-{{ .Name }}"></td> */}
function addUnique(td, name) {
  var input = document.createElement("input");

  input.setAttribute("type", "checkbox");
  input.setAttribute("name", "unique-{{ .Name }}".replace("{{ .Name }}", name))

  td.appendChild(input);
}

{/* <td> */}
{/* <select name="regex-{{ .Name }}"> */}
{/* {{ .Options }} */}
{/* </select> */}
{/* </td> */}
function addData(td, name) {
  var select = document.createElement("select");

  select.setAttribute("name", "regex-{{ .Name }}".replace("{{ .Name }}", name));
  select.innerHTML = document.getElementById("template").value

  td.appendChild(select);
}

{/* <td><input type="checkbox" name="custom-{{ .Name }}"></td> */}
function addCustom(td, name) {
  var input = document.createElement("input");

  input.setAttribute("type", "checkbox");
  input.setAttribute("name", "custom-{{ .Name }}".replace("{{ .Name }}", name));

  td.appendChild(input);
}

{/* <td><input type="text" name="custom-regex-{{ .Name }}"></td> */}
function addREGEX(td, name) {
  var input = document.createElement("input");

  input.setAttribute("type", "text");
  input.setAttribute("name", "custom-regex-{{ .Name }}".replace("{{ .Name }}", name));

  td.appendChild(input);
}

{/* <button type="submit" formaction="/api/v1/generate" target="_blank"> */}
{/* Generate */}
{/* </button> */}
function addDelete(td, id) {
  var button = document.createElement("button");

  button.setAttribute("type", "button");
  button.setAttribute("onclick", 'deleteRow("deletable-{{ .Id }}")'.replace("{{ .Id }}", id));
  button.innerHTML = "Delete"

  td.appendChild(button);
}
