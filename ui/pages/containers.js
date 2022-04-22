import Layout from "../components/layout"

export default function Containers() {

  //...

  return (
    <Layout>
      <p className="help">
        This table displays a list of Docker containers currently being monitored by i5.
        Containers in the running state are proxied unless maintenance mode is active.
      </p>
    </Layout>
  )
}
