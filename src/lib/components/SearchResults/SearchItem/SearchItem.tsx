import React from "react";
import { Topic } from "../../../stores/SemanticSearchStore/interface";

const SearchItem: React.FC<Topic> = ({ context, title }) => {
  return (
    <div className="bg-white rounded-lg shadow-sm hover:shadow transition-shadow duration-200 p-6 mb-4">
      <h3 className="text-xl font-medium text-gray-900 mb-3">{title}</h3>
      <div className="text-gray-600 whitespace-pre-wrap leading-relaxed">
        {context}
      </div>
    </div>
  );
};

export default SearchItem