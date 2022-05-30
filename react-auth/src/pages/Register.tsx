import React, {SyntheticEvent, useState} from 'react';
import { Navigate } from 'react-router-dom';

const Register = () => {

    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [navigate, setNavigate] = useState(false);

    const submit = async (e: SyntheticEvent) => {
        e.preventDefault();

        await fetch('http://localhost:8000/api/register', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({
                name,
                email,
                password
            })
        });

        setNavigate(true)

    }

    if (navigate) {
        return < Navigate to = "/login" />;
    }

    return (

    /*
        <form onSubmit={submit}>
            <h1 className="h3 mb-3 fw-normal">Please register</h1>

            <input className="form-control" placeholder="Name" required
                   onChange={e => setName(e.target.value)}
            />

            <input type="email" className="form-control" placeholder="Email address" required
                   onChange={e => setEmail(e.target.value)}
            />

            <input type="password" className="form-control" placeholder="Password" required
                   onChange={e => setPassword(e.target.value)}
            />

            <button className="w-100 btn btn-lg btn-primary" type="submit">Submit</button>
        </form>

        */


    <body>
        <div className="main">
          <section className="signup">
              <div className="container">
                  <div className="signup-content">
                      <form id="signup-form" className="signup-form" onSubmit={submit}>
                          <h2 className="form-title">Create account</h2>
                          <div className="form-group">
                              <input type="text" className="form-input" name="name" id="name" placeholder="Your Name" onChange={e => setName(e.target.value)} />
                          </div>
                          <div className="form-group">
                              <input type="email" className="form-input" name="email" id="email" placeholder="Your Email" onChange={e => setEmail(e.target.value)}/>
                          </div>
                          <div className="form-group">
                              <input type="text" className="form-input" name="password" id="password" placeholder="Password" onChange={e => setPassword(e.target.value)} />
                              <span className="zmdi zmdi-eye field-icon toggle-password"></span>
                          </div>
                          <div className="form-group">
                              <input type="submit" name="submit" id="submit" className="form-submit" value="Register"/>
                          </div>
                      </form>
                      <p className="loginhere">
                          Have already an account ? <a href="#" className="loginhere-link">Login here</a>
                      </p>
                  </div>
              </div>
          </section>
        </div>
    </body>
    );
};

export default Register;