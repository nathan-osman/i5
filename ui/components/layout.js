import Head from 'next/head'
import Header from './header'

export default function Layout() {
  return (
    <>
      <Head>
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Header />
    </>
  )
}
