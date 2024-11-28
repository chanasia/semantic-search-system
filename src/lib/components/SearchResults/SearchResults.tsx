import React from "react";
import type { Topic } from "../../stores/SemanticSearchStore/interface";
import SearchItem from "./SearchItem/SearchItem";

interface Props {
  results: Topic[]
}

const SearchResult: React.FC<Props> = ({ results }) => {
  return (
    <div className="space-y-4">
      {results.map((result, index) => (
        <SearchItem
          key={index}
          title={result.title.trim()}
          context={result.context.trim()}
        />
      ))}
    </div>
  )
}

export default SearchResult