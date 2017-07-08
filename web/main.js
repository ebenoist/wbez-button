(function() {
  var playing = false;
  var play = function() {
    if (playing) {
      fetch('/stop').then(function() {
        playing = false;
      });
    }

    if (!playing) {
      fetch('/play').then(function() {
        playing = true;
      });
    }
  }

  window.WBEZ = {
    play: play,
  };
})(window)
