namespace DevSmtp.Core.Queries
{
    public class GetMessagesException : Exception
    {
        public GetMessagesException(string message)
            : base(message)
        {
        }

        public GetMessagesException(string message, Exception innerException)
            : base(message, innerException)
        {
        }
    }
}
