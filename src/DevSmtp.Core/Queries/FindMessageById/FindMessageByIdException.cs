namespace DevSmtp.Core.Queries
{
    public class FindMessageByIdException : Exception
    {
        public FindMessageByIdException(string message)
            : base(message)
        {
        }

        public FindMessageByIdException(string message, Exception innerException)
            : base(message, innerException)
        {
        }
    }
}
