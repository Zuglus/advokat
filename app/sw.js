const cacheName = 'v1';
const cacheKeeplist = ['v2'];

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(cacheName).then((cache) => {
      return cache.addAll([
        'index.html',
        'contacts.html',
        'services.html',
        'images/icons/icon-128x128.png',
        'images/icons/icon-192x192.png',
        'images/location.svg',
        'images/logo.svg',
        'images/telegram.svg',
        'images/train.svg',
        'images/viber.svg',
        'images/whatsapp.svg',
        'images/external-link.svg',
        'images/phone.svg',
        'images/email.svg',
        'css/style.css',
        'js/app.js'
      ]);
    })
  );
});

self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((keyList) => {
      return Promise.all(keyList.map((key) => {
        if (cacheKeeplist.indexOf(key) === -1) {
          return caches.delete(key);
        }
      }));
    })
  );
});

self.addEventListener('fetch', (event) => {
  event.respondWith(
    caches.match(event.request).then((resp) => {
      return resp || fetch(event.request).then(async (response) => {
        const cache = await caches.open(cacheName);
        cache.put(event.request, response.clone());
        return response;
      });
    })
  );
});