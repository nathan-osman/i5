import $ from 'jquery';
import axios from "axios";
import { useEffect, useRef, useState } from "react";

/**
 * Form component that makes an AJAX request for data on submission
 */
const Form = ({ url, onData, ...props }) => {

  const formElement = useRef(null);

  const [isLoading, setIsLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState(null);

  useEffect(() => {
    $('input,button', formElement.current).prop('disabled', isLoading);
  }, [isLoading]);

  function handleSubmit(e) {
    e.preventDefault();
    setIsLoading(true);
    setErrorMessage(null);
    const data = $(formElement.current).serializeArray().reduce(
      (m, e) => ({ ...m, [e.name]: e.value }),
      {}
    );
    axios.post(url, data)
      .then(d => {
        onData(d);
      })
      .catch(e => {
        if (e.response && 'error' in e.response.data) {
          setErrorMessage(e.response.data.error);
        } else {
          setErrorMessage("unable to login at this time");
        }
      })
      .finally(() => {
        setIsLoading(false);
      });
  }

  return (
    <form ref={formElement} onSubmit={handleSubmit}>
      {errorMessage ?
        <p className="text-danger">
          <strong>E:</strong>{' '}
          {errorMessage}
        </p> : null
      }
      {props.children}
      <button
        type="submit"
        className="btn btn-primary">
        Submit
      </button>
    </form>
  );
};

export default Form;
