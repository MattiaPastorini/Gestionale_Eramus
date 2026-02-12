import Image from "next/image";
import "bootstrap/dist/css/bootstrap.min.css";
import Link from "next/link";

export default function Home() {
  return (
    <div>
      <h1 className=" d-flex justify-content-center mt-5">Pagina Login</h1>

      {/* Form per il login  */}

      <form className="m-5 border border-3 border-black rounded rounded-3 p-5">
        {/* Input Username */}

        <div className="mb-3">
          <label for="exampleInputEmail1" className="form-label">
            Username
          </label>
          <input
            type="email"
            className="form-control"
            id="exampleInputEmail1"
            aria-describedby="emailHelp"
          ></input>
        </div>

        {/* Input Email */}

        <div className="mb-3">
          <label for="exampleInputEmail1" className="form-label">
            Email address
          </label>
          <input
            type="email"
            className="form-control"
            id="exampleInputEmail1"
            aria-describedby="emailHelp"
          ></input>
        </div>

        {/* Input Password */}

        <div className="mb-3">
          <label for="exampleInputPassword1" className="form-label">
            Password
          </label>
          <input
            type="password"
            className="form-control"
            id="exampleInputPassword1"
          ></input>
        </div>

        {/* Link cambio pagina da Login a Registrazione */}

        <div className="d-flex justify-content-center mb-2">
          <Link href="/register">Non hai un account? Registrati qui.</Link>
        </div>

        <div className="d-flex justify-content-center">
          <button
            type="submit"
            className="btn btn-primary py-2 px-5 bg-danger border-0"
          >
            Accedi
          </button>
        </div>
      </form>
    </div>
  );
}
