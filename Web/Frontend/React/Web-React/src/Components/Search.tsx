import React, { useState } from "react";
import Result_List from "./Result_List";

const bool = "False";

function Search() {
  const [searchbar, getuserinput] = useState("");
  const [result, setResult] = useState(null);
  const submit = async (event) => {
    event.preventDefault();
    const msg = await fetch("http://127.0.0.1:8000/get_query", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ searchbar }),
    });
    const data = await msg.json();
    const bool = "True";
    setResult(data);
  };

  return (
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
              onChange={(s) => getuserinput(s.target.value)}
            ></input>
            <input type="submit" value="submit"></input>
          </form>
        </div>
      </div>

      {result && (
        <div>
          <Result_List data={JSON.stringify(result)} check={bool} />
        </div>
      )}
    </>
  );
}

export default Search;
