import apiClient from '../api/client';

function urlB64ToUint8Array(base64String: string) {
  const padding = '='.repeat((4 - base64String.length % 4) % 4);
  const base64 = (base64String + padding)
    .replace(/\-/g, '+')
    .replace(/_/g, '/');

  const rawData = window.atob(base64);
  const outputArray = new Uint8Array(rawData.length);

  for (let i = 0; i < rawData.length; ++i) {
    outputArray[i] = rawData.charCodeAt(i);
  }
  return outputArray;
}

export async function subscribeToPushNotifications() {
  if (!('serviceWorker' in navigator) || !('PushManager' in window)) {
    console.warn('Push messaging is not supported by this browser.');
    return;
  }

  try {
    // 1. Request Permission
    const permission = await Notification.requestPermission();
    if (permission !== 'granted') {
      console.warn('Permission for notifications was denied');
      return;
    }

    // 2. Register Service Worker
    const registration = await navigator.serviceWorker.register('/sw.js');
    console.log('Service Worker registered with scope:', registration.scope);

    // 3. Wait until Service Worker is ready
    await navigator.serviceWorker.ready;

    // 4. Subscribe to PushManager
    const vapidKey = import.meta.env.VITE_VAPID_PUBLIC_KEY;
    if (!vapidKey) {
      alert("ERROR: VITE_VAPID_PUBLIC_KEY belum terbaca oleh frontend. Pastikan Anda sudah me-restart terminal 'npm run dev'.");
      return;
    }
    const applicationServerKey = urlB64ToUint8Array(vapidKey);
    
    // Check if we already have a subscription
    let subscription = await registration.pushManager.getSubscription();
    if (!subscription) {
      subscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: applicationServerKey
      });
      console.log('New push subscription created.');
    } else {
      console.log('Existing push subscription found.');
    }

    // 5. Send Subscription to Backend
    await apiClient.post('/notifications/subscribe', subscription);
    console.log('Push subscription sent to backend successfully.');
  } catch (error) {
    console.error('Failed to subscribe to push notifications:', error);
    alert('Gagal berlangganan notifikasi: ' + (error as Error).message);
  }
}
