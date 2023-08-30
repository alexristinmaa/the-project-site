import style from "./postListItem.module.css"

import {
    Link
} from "react-router-dom";

export default function PostListItem({post}) {
    function makeNormal(string) {
        string = string.split("-").join(" ");
        return string.charAt(0).toUpperCase() + string.slice(1);
    } 

    return(
        <>
        <article className={style.post}>
            <Link to={"/posts/" + post.Id}>
            <div>
                <img src="/test-thumbnail.jpeg" alt="A Test Thumbail"/>
            </div>
            <div className={style.text}>
                <p className={style.date}>{new Date(post.Date).toLocaleDateString('en-gb', {year:"numeric", month:"long", day:"numeric"}) }</p>
                <h2 className={style.header}>{post.Title}</h2>
                <p className={style.description}>{post.Description}</p>
            </div>
            <div className={style.footer}>
                <div className={style.divisor}></div>
                <Link to={"/posts/page/1?tags=" + post.Tags}><div className={style.tag}>{makeNormal(post.Tags[0])}</div></Link>
            </div>
            </Link>
        </article>
        </>
    )
}