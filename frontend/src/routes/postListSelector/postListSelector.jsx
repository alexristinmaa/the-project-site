import style from "./postListSelector.module.css";

import {
    NavLink
} from "react-router-dom";

export default function PostListSelector({length, current}) {
    let linkList = [current];

    // start from the current letter
    // work outwards with steps until 
    // 4 new page-links have been found
    // Unless dead-point was found on both sides
    // There will be dots added for all
    // non-dead points

    let added = 0;
    let upI = current + 1;
    let downI = current - 1;

    while(added < 4) {
        let oldAdded = added;
        if(upI <= length) {
            linkList.push(upI);
            added++;
        }
        if(downI >= 1) {
            linkList.unshift(downI);
            added++;
        }

        if(oldAdded == added) {
            break;
        }


        upI++;
        downI--;
    }

    // Do we add the ... and the 1 or length
    // Check if our first element is 1
    // Or if our last element is length
    if(linkList[0] != 1) {
        linkList.unshift(0);
        linkList.unshift(1);
    }

    if(linkList[linkList.length-1] != length) {
        linkList.push(0);
        linkList.push(length);
    }

    linkList = linkList.map((v) => {
        if(v == 0) return <span>...</span>;
        if(v == 1) return <NavLink to={`/`} className={({ isActive, isPending }) => isPending ? "pending" : isActive ? style.active : ""} end>{v}</NavLink>
        
        return <NavLink to={`/posts/page/` + v} className={({ isActive, isPending }) => isPending ? "pending" : isActive ? style.active : ""} end>{v}</NavLink>
    })

    return(
        <>
        <div id={style.selector}>
            {
                linkList
            }
        </div>
        </>
    )
}