import { useState } from 'react';
import { Input, Spin, Card, Tag, Alert, Empty, Image } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import { SearchResponse, Topic } from '../types';

const { Search } = Input;

export default function TopicSearchPage() {
  const [searchTerm, setSearchTerm] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [topics, setTopics] = useState<Topic[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSearch = async (value: string) => {
    if (!value.trim()) {
      setTopics([]);
      return;
    }

    setSearchQuery(value);
    setLoading(true);
    setError(null);

    try {
      const backendUrl = import.meta.env.VITE_BACKEND_URL || '';
      const response = await fetch(
        `${backendUrl}/api/v1/topics/search?search_text=${encodeURIComponent(
          value
        )}&limit=10`
      );

      if (!response.ok) {
        throw new Error('Failed to fetch topics');
      }

      const data: SearchResponse = await response.json();
      setTopics(data.results);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
      setTopics([]);
    } finally {
      setLoading(false);
    }
  };

  // Determine if we should show the centered layout
  const showCenteredLayout = !searchQuery && !topics.length;

  return (
    <div className="min-h-screen bg-gray-50">
      {showCenteredLayout ? (
        // Google-style centered layout
        <div className="flex flex-col items-center justify-center min-h-screen px-4">
          <div className="text-center">
            <h1 className="text-7xl font-bold bg-gradient-to-r from-blue-500 to-teal-400 bg-clip-text text-transparent drop-shadow-lg">
              Topic Search
            </h1>
          </div>
          
          <div className="w-full max-w-2xl">
            <div className="search-wrapper">
              <Search
                placeholder="Search topics..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                onSearch={handleSearch}
                className="search-input-large"
                prefix={<SearchOutlined className="search-icon" />}
                allowClear
                size="large"
              />
            </div>
            <style>{`
              .search-wrapper {
                padding: 20px;
              }
              
              .search-input-large {
                width: 100%;
              }
              
              .search-input-large .ant-input-group {
                display: flex;
                border-radius: 24px;
                box-shadow: 0 1px 6px rgba(32, 33, 36, 0.28);
                border: 1px solid transparent;
                overflow: hidden;
              }

              .search-input-large .ant-input-group:hover {
                box-shadow: 0 1px 8px rgba(32, 33, 36, 0.4);
              }

              .search-input-large .ant-input-affix-wrapper {
                flex: 1;
                padding: 12px 16px;
                border: none;
                font-size: 16px;
                background: white;
              }

              .search-input-large .ant-input-affix-wrapper:focus,
              .search-input-large .ant-input-affix-wrapper-focused {
                box-shadow: none;
              }

              .search-input-large .ant-input-affix-wrapper:hover {
                border: none;
              }

              .search-input-large .ant-input {
                font-size: 16px;
                background: transparent;
              }

              .search-input-large .ant-input-group-addon {
                background: transparent;
                border: none;
                padding: 0;
              }

              .search-input-large .ant-input-search-button {
                display: none;
              }

              .search-icon {
                color: #9AA0A6;
                font-size: 20px;
              }

              /* Clear button styling */
              .search-input-large .ant-input-clear-icon {
                color: #9AA0A6;
                font-size: 16px;
              }

              .search-input-large .ant-input-clear-icon:hover {
                color: #202124;
              }
            `}</style>
          </div>
        </div>
      ) : (
        // Regular search results layout
        <div className="py-8">
          <div className="max-w-7xl mx-auto px-4">
            <div className="mb-8">
              <div className="text-center mb-6">
                <a 
                  href="/"
                  className="text-5xl font-bold bg-gradient-to-r from-blue-500 to-teal-400 bg-clip-text text-transparent drop-shadow-lg inline-block hover:opacity-80 transition-opacity cursor-pointer no-underline"
                  onClick={(e) => {
                    e.preventDefault();
                    setSearchTerm('');
                    setSearchQuery('');
                    setTopics([]);
                  }}
                >
                  Topic Search
                </a>
              </div>

              <div className="max-w-2xl mx-auto mb-8">
                <div className="search-wrapper">
                  <Search
                    placeholder="Search topics..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    onSearch={handleSearch}
                    className="search-input-large"
                    prefix={<SearchOutlined className="search-icon" />}
                    allowClear
                    size="large"
                  />
                </div>
              </div>

              {error && (
                <Alert
                  message="Error"
                  description={error}
                  type="error"
                  showIcon
                  className="mb-4"
                />
              )}

              {loading ? (
                <div className="flex justify-center items-center py-8">
                  <Spin size="large" />
                </div>
              ) : topics.length > 0 ? (
                <div className="grid gap-6 grid-cols-1">
                  {topics.map((topic) => (
                    <Card
                      key={topic.id}
                      className="hover:shadow-md transition-shadow overflow-hidden"
                    >
                      <div>
                        <div className="flex items-center justify-between mb-2">
                          <h3 className="text-lg font-medium text-gray-900">{topic.title}</h3>
                          <Tag color="blue">{(topic.similarity * 100).toFixed(1)}% match</Tag>
                        </div>
                    
                        <pre className="text-md text-gray-600 w-full text-wrap">{topic.context}</pre>
                    
                        {topic.tag && (
                          <div className="mb-2">
                            <Tag color="cyan">{topic.tag}</Tag>
                          </div>
                        )}
                    
                        <div className="text-xs text-gray-500 mb-3">
                          Page: {topic.page || 'N/A'}
                        </div>
                    
                        {topic.TopicImages && topic.TopicImages.length > 0 && (
                          <div className="flex gap-4" style={{ maxWidth: '320px' }}>
                            {topic.TopicImages.slice(0, 3).map((image, index) => (
                              <div
                                key={image.id}
                                className="w-96 aspect-square rounded-md overflow-hidden flex-shrink-0"
                              >
                                <div className="w-full h-full bg-gray-50 flex items-center justify-center p-1">
                                  <Image
                                    src={`${import.meta.env.VITE_BACKEND_URL}/api/v1/images/${image.id}`}
                                    alt={`${topic.title} - image ${index + 1}`}
                                    style={{ 
                                      objectFit: 'contain',
                                      maxWidth: '100%',
                                      maxHeight: '100%',
                                      borderRadius: '4px'
                                    }}
                                    preview={{
                                      src: `${import.meta.env.VITE_BACKEND_URL}/api/v1/images/${image.id}`,
                                    }}
                                  />
                                </div>
                              </div>
                            ))}
                          </div>
                        )}
                      </div>
                    </Card>
                  ))}
                </div>
              ) : searchQuery ? (
                <Empty
                  description="No topics found matching your search"
                  className="mt-8"
                />
              ) : null}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}