export interface TopicImage {
  id: number;
  topic_id: string;
  image_path: string;
  created_at: string;
  updated_at: string;
}

export interface Topic {
  id: number;
  title: string;
  context: string;
  page: string;
  tag: string;
  created_at: string;
  updated_at: string;
  TopicImages: TopicImage[];
  similarity: number;
}

export interface SearchResponse {
  results: Topic[];
}