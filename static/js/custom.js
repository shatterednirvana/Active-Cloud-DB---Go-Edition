
$(document).ready(function(){
  $("#queryData").tablesorter({ sortList: [[1,0]] });

  $.ajax({
    type: 'GET',
    async: true,
    url: '/query',
    success: function(data) {
      jsonData = JSON.parse(data);
      for (var key in jsonData) {
        realKey = key.replace('/Entity,', '')
        $('#queryData').append(getRowHTML(realKey, jsonData[key]));
      }
    }
  });

  $("#put").click(function(){
    // TODO: verify that key is not empty
    key = $('#key').val();

    // TODO: verify that val is not empty
    val = $('#val').val();

    $.ajax({
      type: 'POST',
      async: true,
      url: '/put?key=' + key + '&val=' + val,
      success: function(data) {
        console.log(data);
        // TODO: if error is empty, display a success flash
        // TODO: if error is not empty, display it in the error flash
        $('#queryData').append(getRowHTML(key, val));
      }
    });
  });

});

function getRowHTML(key, val) {
  var newRow = "<tr id='" + key + "'><td>" + key + "</td>";
  newRow += "<td>" + val + "</td>";
  newRow += "<td><a href='#' onclick=deleteKey('" + key;
  newRow += "') class='btn danger'>Delete</a></td></tr>";
  return newRow;
}

function deleteKey(key) {
  console.log("key to delete is " + key);

  $.ajax({
    type: 'POST',
    async: true,
    url: '/delete?key=' + key,
    success: function(data) {
      console.log(data);
      $('#' + key).remove();
    }
  });
}
