window.onload = function() {
    const form = document.querySelector('form');
    const statusElement = document.querySelector('#status');
    const messageOutput = document.querySelector('#message-output');
    const openSocketElement = document.querySelector('#open-socket');
    const closeSocketElement = document.querySelector('#close-socket');
    const socketURLInput = document.querySelector('#socket-url-input');

    var socket = null;

    statusElement.textContent = 'Closed';
    statusElement.className = 'tag is-danger';

    form.addEventListener('submit', event => {
        event.preventDefault();
        const messageInput = document.querySelector('#message-input');
        socket.send(messageInput.value);
        messageInput.value = '';
    });

    openSocketElement.onclick = function() {
        socket = new WebSocket(socketURLInput.value);

        socket.onopen = function(e) {
            console.log('Socket open:', e);
            statusElement.textContent = 'Opened';
            statusElement.className = 'tag is-success'
        }
    
        socket.onmessage = function(e) {
            console.log('Message:', e);
            let data = messageOutput.innerHTML;
            messageOutput.innerHTML = data.concat(e.data, '\r\n');
            messageOutput.scrollTop = messageOutput.scrollHeight;
        }
    
        socket.onclose = function(e) {
            console.log('Socket closed:', e);
            statusElement.textContent = 'Closed';
            statusElement.className = 'tag is-danger';
        }

        socket.onerror = function(e) {
            console.log(e);
            let data = messageOutput.innerHTML;
            messageOutput.innerHTML = data.concat('Error occurred\r\n');
        }
    }

    closeSocketElement.onclick = function() {
        socket.close();
    }
}
