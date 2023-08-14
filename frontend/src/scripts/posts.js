export async function getPosts() {
    let promise = await fetch("/api/posts");
    let posts = await promise.json();

    return posts;
}