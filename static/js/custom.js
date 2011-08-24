
$(document).ready(function(){
  $("#queryData").tablesorter({ sortList: [[1,0]] });

  $('#success-flash-div').hide();
  $('#error-flash-div').hide();
 
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
    key = $('#key').val();
    if (key == '') {
      $('#error-flash').text('Key cannot be empty.'); 
      showErrorDiv();
      return;
    }

    val = $('#val').val();
    if (val == '') {
      $('#error-flash').text('Value cannot be empty.'); 
      showErrorDiv();
      return;
    }

    $.ajax({
      type: 'POST',
      async: true,
      url: '/put?key=' + key + '&val=' + val,
      success: function(data) {
        jsonData = JSON.parse(data);
        error = jsonData['error']
        if (error == '') {
          $('#success-flash').text('Data saved successfully!');
          showSuccessDiv();
          $('#queryData').append(getRowHTML(key, val));
        } else {
          errorMsg = 'There was a problem saving your data: ' + error;
          $('#error-flash').text(errorMsg);
          showErrorDiv();
        }
      }
    });
  });

});

function showSuccessDiv() {
  $('#success-flash-div').show();
  $('#error-flash-div').hide();
}

function showErrorDiv() {
  $('#success-flash-div').hide();
  $('#error-flash-div').show();
}

function getRowHTML(key, val) {
  var newRow = "<tr id='" + key + "'><td>" + key + "</td>";
  newRow += "<td>" + val + "</td>";
  newRow += "<td><a href='#' onclick=deleteKey('" + key;
  newRow += "') class='btn danger'>Delete</a></td></tr>";
  return newRow;
}

function deleteKey(key) {
  $.ajax({
    type: 'POST',
    async: true,
    url: '/delete?key=' + key,
    success: function(data) {
      console.log(data);
      $('#' + key).remove();
      $('#success-flash').text('"' + key + '" deleted successfully.');
      showSuccessDiv();
    }
  });
}
