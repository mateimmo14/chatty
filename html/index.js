let signedin = false;
let username = "";
let ws = null;
let text = null;
let button = document.getElementById("sendbutton");
let input = document.getElementById("message");

button.addEventListener("click", () => {
    if (!signedin) {
        if (input.value.trim() !== "") {
            username = input.value.trim();
            signedin = true;
            document.getElementById("body").innerHTML = `
                <textarea id="chat" class="item" readonly></textarea>
                <div>
                    <input id="message" class="item" placeholder="Type a message">
                    <button id="sendbutton" class="item">Send</button>
                </div>
            `;
            text = document.getElementById("chat");
            button = document.getElementById("sendbutton");
            input = document.getElementById("message");
            ws = new WebSocket(`ws://${window.location.host}/connect`);
            ws.onmessage = (event) => {
                text.textContent += event.data;
            };
            button.addEventListener("click", () => {
                ws.send(`${username}: ${input.value}`);
                input.value = "";
            });
        }
    }
});