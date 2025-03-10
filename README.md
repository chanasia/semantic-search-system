```bash
docker volume create semantic-search-system_minio_data

docker volume create semantic-search-system_projectdb64

docker run --rm \
  -v semantic-search-system_minio_data:/volume \
  -v "$(pwd)":/backup \
  alpine sh -c "cd /volume && tar xzvf /backup/semantic-search-system_minio_data.tar.gz"


docker run --rm \
  -v semantic-search-system_projectdb64:/volume \
  -v "$(pwd)":/backup \
  alpine sh -c "cd /volume && tar xzvf /backup/semantic-search-system_projectdb64.tar.gz" 

docker-compose up -d

docker-compose down --rmi
```