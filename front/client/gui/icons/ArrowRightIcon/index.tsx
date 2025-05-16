import React from 'react'

const ArrowRight = (props) => (
	<svg className={(props.css || "")} width={(props.width || 800) + "px"} height={(props.height || 800) + "px"} viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
	<path d="M5 12H19M19 12L13 6M19 12L13 18" stroke="#000000" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
	</svg>
)

export default ArrowRight