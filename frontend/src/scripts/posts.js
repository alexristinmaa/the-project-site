/* export async function getPosts() {
    let promise = await fetch("/api/posts");
    let posts = await promise.json();

    return posts;
} */

export async function getFeaturedPost() {
    return {Date:"9 August 2023", Title:"Injecting XXS Payloads into Web Applications using Three.js", Description:"Bypassing specially designed barriers using secret tricks of the trade. Follow me on a research journey into MacOS aliases.", Tag:"hacking", Id:"bypassing-specially-designed-barriers"};
}

export async function getPosts({search, tags, page}) {
    page ||= 1;
    tags ||= [];
    search ||= "";

    let data = {
        page: page,
        tags: tags,
        search: search
    }

    let response = await fetch("/api/posts/", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            // 'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: JSON.stringify(data)
    });

    let jsonD = await response.json();

    console.log(jsonD);

    return jsonD;
}

export async function getTags() {
    let response = await fetch("/api/getTags");

    let jsonD = await response.json();

    return jsonD;
} 