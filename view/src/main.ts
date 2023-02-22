import './style.css'
import typescriptLogo from './typescript.svg'
import { setupCounter } from './counter'


const socket = new WebSocket('ws://localhost:5000/panel');
socket.addEventListener("message", (event) => {
  console.log('Socket messeger: ', event.data)
  console.log(event)
})


// setTimeout(() => {
//   socket.send("entry");
// }, 5000)

// socket.onmessage = ({data}) => {
//   console.log(data)
// };


document.querySelector<HTMLDivElement>('#app')!.innerHTML = `
  <div>
    <a href="https://vitejs.dev" target="_blank">
      <img src="/vite.svg" class="logo" alt="Vite logo" />
    </a>
    <a href="https://www.typescriptlang.org/" target="_blank">
      <img src="${typescriptLogo}" class="logo vanilla" alt="TypeScript logo" />
    </a>
    <h1>Vite + TypeScript</h1>
    <div class="card">
      <button id="counter" type="button"></button>
    </div>
    <p class="read-the-docs">
      Click on the Vite and TypeScript logos to learn more
    </p>
  </div>
`

setupCounter(document.querySelector<HTMLButtonElement>('#counter')!)
