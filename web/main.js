(function(window, document) {
  var status = 'stopped';

  var setStatus = function(body) {
    document.querySelector('#status').innerText = 'WBEZ: '+ body;
    status = body;
  };

  var onLoad = function() {
    fetch('/status').then(function(resp) {
      return resp.text();
    }).then(setStatus);
  }

  var onClick = function() {
    if (status === 'streaming') {
      fetch('/stop').then(function(resp) {
        return resp.text();
      }).then(setStatus);
    }

    if (status === 'stopped') {
      fetch('/play').then(function(resp) {
        return resp.text();
      }).then(setStatus);
    }
  }

  window.WBEZ = {
   onClick: onClick,
   onLoad: onLoad,
  };
})(window, document)
