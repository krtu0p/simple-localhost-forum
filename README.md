Installation
Build simple-localhost-forum from the source and install dependencies:
1. Clone the repository:
git clone https://github.com/username/simple-localhost-forum

2. Navigate to the project directory:
cd simple-localhost-forum

3. Install the dependencies:
Using Docker:
docker build -t simple-localhost-forum .

Using go modules:
go build

4. Run the project with:
Using Docker:
docker run --image simple-localhost-forum
