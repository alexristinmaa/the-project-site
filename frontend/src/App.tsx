import { useEffect, useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

interface Post {
  title : string,
  author : string
}

function App() {
  const [posts, setPosts] = useState([]);

  const makePost = (post : Post) => <h1>{post.title} by {post.author}</h1>;

  const getPosts = async () => {
    let res = await fetch("/posts");
    let jsonData = await res.json();

    setPosts(jsonData);
  }

  useEffect(() => {getPosts()}, []);

  return (
    <>
      <div>
        <a href="https://vitejs.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <div className="card">
        <button onClick={() => alert("Hi")}>
          {"Special!"}
        </button>
        <p>
        {posts.map(makePost)}
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  )
}

export default App
