import ContainerList from '../components/containerlist'

export default function Containers() {

  return (
    <>
      <div className="title">Containers</div>
      <div className="secondary">
        This page lists all containers running on the host that are recognized by i5 (containers with the <kbd>i5.*</kbd> set of tags). You can use the tools available on the right to control the containers.
      </div>
      <ContainerList />
    </>
  )
}
