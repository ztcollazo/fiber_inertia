const Index = ({ name, req, start }: { name: string, req: string, start: string }) => {
  return (
    <div>
      <h1>Hello, {name}!</h1>
      <p>Server started at: {start}</p>
      <p>Page requested at: {req}</p>
    </div>
  )
}

export default Index;
