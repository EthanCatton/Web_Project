function Result_List(data) {
  // gets rid of outer brackets
  console.log(data.data);
  const oneparse = JSON.parse(data.data);
  //gets to the point of an array which each contain an object of {URL:https://example.com,match_score:10}
  const twoparse = JSON.parse(oneparse.result);
  return (
    <>
      <div className="search_container">
        <h2 className="searchbar">Results:</h2>
      </div>
      <ul className="Result_List">
        {twoparse.map((item, index) => (
          <li key={index}>
            <div className="result">
              <a href={item.URL} target="_blank" rel="noopener noreferrer">
                {item.Name}
                <br></br>
              </a>{" "}
              {"     score: "}
              {item.match_score}
            </div>
          </li>
        ))}
      </ul>
    </>
  );
}

export default Result_List;
