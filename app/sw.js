const staticCacheName = 'static-cache-v0';

const staticAssets = [
  '/',
  '/index.html',
  '/contacts.html',
  '/services.html',
  '/images/icons/icon-128x128.png',
  '/images/icons/icon-192x192.png',
  '/images/location.svg',
  '/images/phone.svg',
  '/images/email.svg',
  '/css/style.css',
  '/js/app.js',
  '/js/index.js'
];

self.addEventListener('install', async event => {

  const cache = await caches.open(staticCacheName);
  await cache.addAll(staticAssets);
  console.log('Service worker has been installed');
});

self.addEventListener('activate', async event => {
  const cachesKeys = await caches.keys();
  const checkKeys = cachesKeys.map(async key => {
   if (staticCacheName !== key) {
    await caches.delete(key);
  }
  });
  await Promise.all(checkKeys);
  console.log('Service worker has been activated');
});

self.addEventListener('fetch', event => {
  console.log(`Trying to fetch ${event.request.url}`);
  event.respondWith(caches.match(event.request).then(cachedResponse => {
    return cachedResponse || fetch(event.request)
  }));
});