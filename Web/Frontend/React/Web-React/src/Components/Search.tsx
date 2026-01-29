import React, { useState } from "react";
import Result_List from "./Result_List";

function Search() {
  const [searchbar, getuserinput] = useState("");
  const [result, set_result] = useState(null);
  //on submit press
  const submit = async (event) => {
    // stop refresh
    event.preventDefault();
    // query function
    const msg = await fetch("http://127.0.0.1:8000/get_query", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ searchbar }),
    });
    const data = await msg.json();
    // adjust ui with new elements
    set_result(data);
  };

  return (
    // returns section for searching which will then adapt with results
    <>
      <div className="search_container">
        <div className="searchbar">
          <form onSubmit={submit}>
            <label htmlFor="searchbar">Search for: </label>
            <input
              type="text"
              id="searchbar"
              name="searchbar"
              value={searchbar}
              onChange={(usersearch) => getuserinput(usersearch.target.value)}
            ></input>
            <input type="submit" value="submit"></input>
          </form>
        </div>
      </div>

      {result && (
        <div>
          <Result_List data={JSON.stringify(result)} />
        </div>
      )}
    </>
  );
}

export default Search;
