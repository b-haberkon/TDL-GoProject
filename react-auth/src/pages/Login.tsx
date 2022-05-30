import React, { SyntheticEvent, useState } from 'react';
import { Link, Navigate } from 'react-router-dom';
import "../Css/Login.css";


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
                            <h2 className="form-title">Create account</h2>
                            <div className="form-group">
                                <input type="text" className="form-input" name="name" id="name" placeholder="Your Name"/>
                            </div>
                            <div className="form-group">
                                <input type="email" className="form-input" name="email" id="email" placeholder="Your Email"/>
                            </div>
                            <div className="form-group">
                                <input type="text" className="form-input" name="password" id="password" placeholder="Password"/>
                                <span className="zmdi zmdi-eye field-icon toggle-password"></span>
                            </div>
                            <div className="form-group">
                                <input type="password" className="form-input" name="re_password" id="re_password" placeholder="Repeat your password"/>
                            </div>
                            <div className="form-group">
                                <input type="checkbox" name="agree-term" id="agree-term" className="agree-term" />
                                <label  className="label-agree-term"><span><span></span></span>I agree all statements in  <a href="#" className="term-service">Terms of service</a></label>
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