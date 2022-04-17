import Head from 'next/head'
import Header from './header'

export default function Layout({ children }) {
  return (
    <>
      <Head>
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Header />
      <div className="container">
        {children}
      </div>
    </>
  )
}
