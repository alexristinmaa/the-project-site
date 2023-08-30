import * as React from 'react';
import * as ReactDOM from 'react-dom/client';
import {
  createBrowserRouter,
  RouterProvider,
} from 'react-router-dom';

import Root, { loader as rootLoader } from "./routes/root";
import About from "./routes/about";
import Me from "./routes/me";
import Post, {loader as postLoader} from './routes/post/post'
import SearchPage, {loader as searchLoader} from './routes/searchPage/searchPage';
import SearchResults, {loader as searchResultsLoader} from './routes/searchResults/searchResults';

import ErrorPage from './error-page';

import "./main.css";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    loader: rootLoader,
    errorElement: <ErrorPage />,
    children: [
      {
        path: "",
        element: <SearchResults />,
        loader: searchResultsLoader
      }
    ]
  },
  {
    path: "posts/:postId",
    element: <Post />,
    loader: postLoader
  },
  {
    path: "/me",
    element: <Me />
  },
  {
    path: "/about",
    element: <About />
  },
  {
    path: "posts/page/:pageNumber",
    element: <SearchPage />,
    loader: searchLoader,
    children: [
      {
        path: "",
        element: <SearchResults />,
        loader: searchResultsLoader
      }

    ]
  }
]);

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);