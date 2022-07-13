import React, { SyntheticEvent, useState } from 'react';
import { Link, Navigate } from 'react-router-dom';
import "../css/Login.css";


const Login = (props: { setName: (name: string) => void }) => {

    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [redirect, setRedirect] = useState(false);

    const submit = async (e: SyntheticEvent) => {

        e.preventDefault();

        const response = await fetch('http://localhost:8000/api/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({
                email,
                password
            })
        });

        const content = await response.json();

        setRedirect(true);
        props.setName(content.name);

    }

    if (redirect) {
        return <Navigate to="/home" />
    }

    return ( /*
        <form onSubmit={submit}>
            <h1 className="h3 mb-3 fw-normal">Please sign in</h1>
            <input type="email" className="form-control" placeholder="Email address" required
                onChange={e => setEmail(e.target.value)}
            />

            <input type="password" className="form-control" placeholder="Password" required
                onChange={e => setPassword(e.target.value)}
            />

            <button className="w-100 btn btn-lg btn-primary" type="submit">Sign in</button>
        </form> */
        <body>
          <div className="main">
            <section className="signup">
                <div className="container">
                    <div className="signup-content">
                        <form id="signup-form" className="signup-form" onSubmit={submit}>
                            <h2 className="form-title">Log in</h2>
                            <div className="form-group">
                                <input type="email" className="form-input" name="email" id="email" placeholder="Your Email" onChange={e => setEmail(e.target.value)}/>
                            </div>
                            <div className="form-group">
                                <input type="password" className="form-input" name="password" id="password" placeholder="Password" onChange={e => setPassword(e.target.value)}/>
                                <span className="zmdi zmdi-eye field-icon toggle-password"></span>
                            </div>
                            <div className="form-group">
                                <input type="submit" name="submit" id="submit" className="form-submit" value="Sign up"/>
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
}

export default Login;