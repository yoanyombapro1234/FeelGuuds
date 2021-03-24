import React from 'react'
import {Row, Col,Container, Navbar, Nav, NavDropdown, Form, FormControl, Button} from 'react-bootstrap'

function DropdownOption(element, tagline) {
	return (

		<Container fluid>
			<Row className="py-1">
				<Col sm={4} md={4} lg={4}>

				</Col>

				<Col sm={8} md={8} lg={8}>
					<Row>
						<Nav.Link href="/{element}">{element}</Nav.Link>
					</Row>
					<Row>
						<p>{tagline}</p>
					</Row>
				</Col>
			</Row>
		</Container>
	)
}

function Inspiration() {
	return (
		<NavDropdown title="Inspiration" id="basic-nav-dropdown" className="dropdown">
			<Container fluid>
				<Row>
				<Col sm={8} md={8} lg={8}>
					{DropdownOption("Explore Design Work", "Trending Designs To Inspire You")}
					{DropdownOption("New & Noteworthy", "Up & Coming Noteworthy Items")}
					{DropdownOption("Active & Trending", "Products Customers Are Rifling Over")}
					{DropdownOption("Companies", "Interesting Companies To Follow")}
				</Col>
				<Col sm={4} md={4} lg={4}>
				 	Browse
					<Nav.Link href="#link">Products Of The Day</Nav.Link>
					<Nav.Link href="#link">New</Nav.Link>
					<Nav.Link href="#link">Under $10</Nav.Link>
					<Nav.Link href="#link">Minority Owned</Nav.Link>
					<Nav.Link href="#link">Favorites</Nav.Link>
					<Nav.Link href="#link">Trending</Nav.Link>
				</Col>
				</Row>
			</Container>
		</NavDropdown>
	)
}

const Header = () => {
	return (
		<div>
			<Navbar bg="light" expand="lg">
			<Navbar.Brand href="/home">FeelGuuds</Navbar.Brand>
			<Navbar.Toggle aria-controls="basic-navbar-nav" />
			<Navbar.Collapse id="basic-navbar-nav">
				<Nav className="mr-auto">
				<Nav.Link href="#home">Inspiration</Nav.Link>
				<Nav.Link href="#link">Businesses</Nav.Link>
				{Inspiration()}
				</Nav>
			</Navbar.Collapse>
			<Navbar.Collapse id="basic-navbar-nav">
				<Nav className="ml-auto">
					<Nav.Link href="/cart"> <i className="fas fa-shopping-cart"/> Cart</Nav.Link>
					<Nav.Link href="/login"> <i className="fas fa-user"/> Sign In</Nav.Link>
				</Nav>
			</Navbar.Collapse>
			</Navbar>
		</div>
	)
}

export default Header
