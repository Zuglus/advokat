window.addEventListener('load', () => {

  if ('serviceWorker' in navigator) {
    navigator.serviceWorker.register('../sw.js')
    .then(registration => {
      console.log('Service worker successfully registered', registration);
    })
    .catch(console => {
      console.log('Service worker registration failed', error);
    });
  }
});