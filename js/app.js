document.addEventListener('DOMContentLoaded', function () {
  var burgerBtn = document.querySelector('.burger-btn');
  var mobileMenu = document.querySelector('.mobile-menu');

  if (burgerBtn && mobileMenu) {
    burgerBtn.addEventListener('click', function () {
      var isOpen = mobileMenu.classList.toggle('open');
      burgerBtn.setAttribute('aria-expanded', isOpen ? 'true' : 'false');
      burgerBtn.setAttribute('aria-label', isOpen ? 'Закрыть меню' : 'Открыть меню');
    });
  }
});
