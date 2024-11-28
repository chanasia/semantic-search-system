import { useState } from "react";
import SearchIcon from '@mui/icons-material/Search';
import CloseIcon from '@mui/icons-material/Close';

const SearchText = () => {
  const [value, setValue] = useState("");

  const handleClear = () => {
    setValue("");
  };

  return (
    <div className="mb-6">
      <label className="relative flex items-center">
        <div className="absolute left-4 text-gray-400">
          <SearchIcon />
        </div>

        <input
          type="search"
          value={value}
          onChange={(e) => setValue(e.target.value)}
          placeholder="ค้นหา..."
          className="w-full py-3 pl-12 pr-10 bg-white border rounded-lg shadow-sm 
          placeholder-gray-400 text-gray-700
          focus:outline-none focus:ring-2 focus:ring-blue-100 focus:border-blue-400
          transition-all duration-200"
        />

        {value && (
          <button
            onClick={handleClear}
            className="absolute right-4 text-gray-400 hover:text-gray-600 transition-colors"
          >
            <CloseIcon />
          </button>
        )}
      </label>
    </div>
  );
};

export default SearchText;