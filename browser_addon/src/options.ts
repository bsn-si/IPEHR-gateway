import '../styles/options.scss';

const saveOptions = () => {
  const url = (document.getElementById('url') as HTMLInputElement).value;
  const username = (document.getElementById('username') as HTMLInputElement).value;
  const password = (document.getElementById('password') as HTMLInputElement).value;
  const systemId = (document.getElementById('systemId') as HTMLInputElement).value;
  const token = (document.getElementById('token') as HTMLInputElement).value;

  chrome.storage.sync.set(
    { url, username, password, token, systemId },
    () => {
      const status = document.getElementById('status');
      
      status.textContent = 'Options saved.';
      setTimeout(() => {
        status.textContent = '';
      }, 750);
    }
  );
};

const restoreOptions = () => {
  chrome.storage.sync.get(
    ['url', 'username', 'password', 'token', 'systemId' ],
    (items) => {
      console.log('restoreOptions', items);
      (document.getElementById('url') as HTMLInputElement).value = items.url;
      (document.getElementById('username') as HTMLInputElement).value = items.username;
      (document.getElementById('password') as HTMLInputElement).value = items.password;
      (document.getElementById('token') as HTMLInputElement).value = items.token;
      (document.getElementById('systemId') as HTMLInputElement).value = items.systemId;
    }
  );
};

document.addEventListener('DOMContentLoaded', restoreOptions);
document.getElementById('save').addEventListener('click', saveOptions);
