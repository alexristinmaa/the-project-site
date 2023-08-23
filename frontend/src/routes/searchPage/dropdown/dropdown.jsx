import style from "./dropdown.module.css";
 
export default function Dropdown({tags, activeState, changeActive}) {
    function toggle(list, val) {
        if(list.includes(val)) {
            let index = list.indexOf(val);
            return list.slice(0, index).concat(list.slice(index+1,list.length));
        }

        return [...list, val];
    }

    function changeTag(e) {
         let name = e.target.getAttribute("name");

         changeActive(toggle(activeState, name));
    }

    function makeNormal(tag) {
        tag = tag.split("-").join(" ");
        return tag.charAt(0).toUpperCase() + tag.slice(1);
    }

    let summary = <summary value={activeState} id={style.active}><b>{activeState.length}</b> active tags</summary>;

    if(activeState.length == 1) {
        summary = <summary value={activeState} id={style.active}>By tag: <b>{makeNormal(activeState[0])}</b></summary>
    } else if(activeState.length == 0) {
        summary = <summary value={activeState} id={style.active}>By tag: <b>All tags</b></summary>
    }

    console.log("dropdown", tags)

    return(
        <>
            <details id={style.dropdown}>
                {summary}

                <div id={style.droplist}>
                    {
                        tags.map((tag) => {
                            let classString = `${style.option} ${activeState.includes(tag.Name) ? style.selected : ""}`;
                            return <div name={tag.Name} className={classString} onClick={changeTag}>{makeNormal(tag.Name)} - {tag.Count}</div>;
                        })
                    }
                </div>

            </details>
        </>
    )
}