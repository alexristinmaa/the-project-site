import { useEffect, useState } from 'react'

import PostItem from './post/post.tsx'

interface PostObj {
  Title : string,
  Author : string,
  Description : string
}

function PostList() {
  const [posts, setPosts] = useState<PostObj[]>([]);

  const getPosts = async () => {
    let res = await fetch("/posts");
    let jsonData : PostObj[] = await res.json();

    setPosts(jsonData);
  }

  useEffect(() => {getPosts()}, []);

  return (
    <div className="posts">
        {posts.map((post : PostObj) => <PostItem name={post.Title} author={post.Author} description={post.Description} />)}
    </div>
  )
}

export default PostList
