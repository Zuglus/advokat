document.addEventListener('DOMContentLoaded', function () {
  var buttons = document.querySelectorAll('.tabs-btn');
  var contents = document.querySelectorAll('.tabs-content');

  buttons.forEach(function (btn) {
    btn.addEventListener('click', function () {
      var tabName = this.getAttribute('data-tab');

      buttons.forEach(function (b) { b.classList.remove('active'); });
      contents.forEach(function (c) { c.classList.remove('active'); });

      this.classList.add('active');
      var target = document.querySelector('[data-tab-content="' + tabName + '"]');
      if (target) target.classList.add('active');
    });
  });
});
