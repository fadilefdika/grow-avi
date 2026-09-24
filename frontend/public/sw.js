self.addEventListener('push', function(event) {
  try {
    let title = 'Notifikasi GROW AVI';
    let options = {
      body: 'Anda mendapat pesan baru',
      data: { url: '/' }
    };

    if (event.data) {
      const data = event.data.json();
      title = data.title || title;
      options.body = data.body || options.body;
      options.data.url = data.url || options.data.url;
    }

    event.waitUntil(
      self.registration.showNotification(title, options)
    );
  } catch (error) {
    console.error('Push event error:', error);
    // Fallback safe notification if parsing fails
    event.waitUntil(
      self.registration.showNotification('Pesan Baru', {
        body: 'Anda mendapat notifikasi baru (Fallback)'
      })
    );
  }
});

self.addEventListener('notificationclick', function(event) {
  event.notification.close();
  event.waitUntil(
    clients.openWindow(event.notification.data.url || '/')
  );
});
