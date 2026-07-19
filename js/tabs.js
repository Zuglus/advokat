document.addEventListener('DOMContentLoaded', function () {
  var buttons = document.querySelectorAll('.tabs-btn');
  var contents = document.querySelectorAll('.tabs-content');

  buttons.forEach(function (btn) {
    btn.addEventListener('click', function () {
      var tabName = this.getAttribute('data-tab');

      buttons.forEach(function (b) {
        b.classList.remove('active');
        b.setAttribute('aria-selected', 'false');
      });
      contents.forEach(function (c) { c.classList.remove('active'); });

      this.classList.add('active');
      this.setAttribute('aria-selected', 'true');
      var target = document.querySelector('[data-tab-content="' + tabName + '"]');
      if (target) target.classList.add('active');
    });

    btn.addEventListener('keydown', function (e) {
      var index = Array.prototype.indexOf.call(buttons, this);
      var nextIndex = -1;

      if (e.key === 'ArrowRight' || e.key === 'ArrowDown') {
        nextIndex = (index + 1) % buttons.length;
      } else if (e.key === 'ArrowLeft' || e.key === 'ArrowUp') {
        nextIndex = (index - 1 + buttons.length) % buttons.length;
      } else if (e.key === 'Home') {
        nextIndex = 0;
      } else if (e.key === 'End') {
        nextIndex = buttons.length - 1;
      }

      if (nextIndex >= 0) {
        e.preventDefault();
        buttons[nextIndex].focus();
        buttons[nextIndex].click();
      }
    });
  });
});
