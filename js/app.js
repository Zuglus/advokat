document.addEventListener('DOMContentLoaded', function () {
  var burgerBtn = document.querySelector('.burger-btn');
  var mobileMenu = document.querySelector('.mobile-menu');

  if (burgerBtn && mobileMenu) {
    burgerBtn.addEventListener('click', function () {
      mobileMenu.classList.toggle('open');
    });
  }
});
