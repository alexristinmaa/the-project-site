interface PostData {
    name : string,
    author : string,
    description : string
}

function PostItem({name, author, description} : PostData) {
  return (
    <>
    <div className="post">
        <div className="postTitle"><span>{name}</span> by <span>{author}</span></div>
        <div className="postDescription">{description}</div>
    </div>
    </>
  )
}

export default PostItem
